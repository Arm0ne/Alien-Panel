import{Ci as e,Cr as t,Di as n,Ei as r,Fr as i,Oi as a,Pr as o,Si as s,Sr as c,Ti as l,ca as u,ea as d,ei as f,or as p,sa as m,ui as h,wi as g,xi as _}from"./router-Bly_k9Sk.js";import{a as v}from"./index-CRVuJ7Ja.js";function y(e,t=`default`,n=[]){let{children:r}=e;if(typeof r==`object`&&r&&!Array.isArray(r)){let e=r[t];if(typeof e==`function`)return e()}return n}var b=_([s(`descriptions`,{fontSize:`var(--n-font-size)`},[s(`descriptions-separator`,`
 display: inline-block;
 margin: 0 8px 0 2px;
 `),s(`descriptions-table-wrapper`,[s(`descriptions-table`,[s(`descriptions-table-row`,[s(`descriptions-table-header`,{padding:`var(--n-th-padding)`}),s(`descriptions-table-content`,{padding:`var(--n-td-padding)`})])])]),l(`bordered`,[s(`descriptions-table-wrapper`,[s(`descriptions-table`,[s(`descriptions-table-row`,[_(`&:last-child`,[s(`descriptions-table-content`,{paddingBottom:0})])])])])]),g(`left-label-placement`,[s(`descriptions-table-content`,[_(`> *`,{verticalAlign:`top`})])]),g(`left-label-align`,[_(`th`,{textAlign:`left`})]),g(`center-label-align`,[_(`th`,{textAlign:`center`})]),g(`right-label-align`,[_(`th`,{textAlign:`right`})]),g(`bordered`,[s(`descriptions-table-wrapper`,`
 border-radius: var(--n-border-radius);
 overflow: hidden;
 background: var(--n-merged-td-color);
 border: 1px solid var(--n-merged-border-color);
 `,[s(`descriptions-table`,[s(`descriptions-table-row`,[_(`&:not(:last-child)`,[s(`descriptions-table-content`,{borderBottom:`1px solid var(--n-merged-border-color)`}),s(`descriptions-table-header`,{borderBottom:`1px solid var(--n-merged-border-color)`})]),s(`descriptions-table-header`,`
 font-weight: 400;
 background-clip: padding-box;
 background-color: var(--n-merged-th-color);
 `,[_(`&:not(:last-child)`,{borderRight:`1px solid var(--n-merged-border-color)`})]),s(`descriptions-table-content`,[_(`&:not(:last-child)`,{borderRight:`1px solid var(--n-merged-border-color)`})])])])])]),s(`descriptions-header`,`
 font-weight: var(--n-th-font-weight);
 font-size: 18px;
 transition: color .3s var(--n-bezier);
 line-height: var(--n-line-height);
 margin-bottom: 16px;
 color: var(--n-title-text-color);
 `),s(`descriptions-table-wrapper`,`
 transition:
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `,[s(`descriptions-table`,`
 width: 100%;
 border-collapse: separate;
 border-spacing: 0;
 box-sizing: border-box;
 `,[s(`descriptions-table-row`,`
 box-sizing: border-box;
 transition: border-color .3s var(--n-bezier);
 `,[s(`descriptions-table-header`,`
 font-weight: var(--n-th-font-weight);
 line-height: var(--n-line-height);
 display: table-cell;
 box-sizing: border-box;
 color: var(--n-th-text-color);
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `),s(`descriptions-table-content`,`
 vertical-align: top;
 line-height: var(--n-line-height);
 display: table-cell;
 box-sizing: border-box;
 color: var(--n-td-text-color);
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `,[e(`content`,`
 transition: color .3s var(--n-bezier);
 display: inline-block;
 color: var(--n-td-text-color);
 `)]),e(`label`,`
 font-weight: var(--n-th-font-weight);
 transition: color .3s var(--n-bezier);
 display: inline-block;
 margin-right: 14px;
 color: var(--n-th-text-color);
 `)])])])]),s(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color);
 --n-merged-td-color: var(--n-td-color);
 --n-merged-border-color: var(--n-border-color);
 `),n(s(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color-modal);
 --n-merged-td-color: var(--n-td-color-modal);
 --n-merged-border-color: var(--n-border-color-modal);
 `)),a(s(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color-popover);
 --n-merged-td-color: var(--n-td-color-popover);
 --n-merged-border-color: var(--n-border-color-popover);
 `))]),x=`DESCRIPTION_ITEM_FLAG`;function S(e){return typeof e==`object`&&e&&!Array.isArray(e)?e.type&&e.type.DESCRIPTION_ITEM_FLAG:!1}var C=m({name:`Descriptions`,props:Object.assign(Object.assign({},p.props),{title:String,column:{type:Number,default:3},columns:Number,labelPlacement:{type:String,default:`top`},labelAlign:{type:String,default:`left`},separator:{type:String,default:`:`},size:String,bordered:Boolean,labelClass:String,labelStyle:[Object,String],contentClass:String,contentStyle:[Object,String]}),slots:Object,setup(e){let{mergedClsPrefixRef:n,inlineThemeDisabled:i,mergedComponentPropsRef:a}=t(e),o=d(()=>e.size||a?.value?.Descriptions?.size||`medium`),s=p(`Descriptions`,`-descriptions`,b,v,e,n),l=d(()=>{let{bordered:t}=e,n=o.value,{common:{cubicBezierEaseInOut:i},self:{titleTextColor:a,thColor:c,thColorModal:l,thColorPopover:u,thTextColor:d,thFontWeight:f,tdTextColor:p,tdColor:m,tdColorModal:h,tdColorPopover:g,borderColor:_,borderColorModal:v,borderColorPopover:y,borderRadius:b,lineHeight:x,[r(`fontSize`,n)]:S,[r(t?`thPaddingBordered`:`thPadding`,n)]:C,[r(t?`tdPaddingBordered`:`tdPadding`,n)]:w}}=s.value;return{"--n-title-text-color":a,"--n-th-padding":C,"--n-td-padding":w,"--n-font-size":S,"--n-bezier":i,"--n-th-font-weight":f,"--n-line-height":x,"--n-th-text-color":d,"--n-td-text-color":p,"--n-th-color":c,"--n-th-color-modal":l,"--n-th-color-popover":u,"--n-td-color":m,"--n-td-color-modal":h,"--n-td-color-popover":g,"--n-border-radius":b,"--n-border-color":_,"--n-border-color-modal":v,"--n-border-color-popover":y}}),u=i?c(`descriptions`,d(()=>{let t=``,{bordered:n}=e;return n&&(t+=`a`),t+=o.value[0],t}),l,e):void 0;return{mergedClsPrefix:n,cssVars:i?void 0:l,themeClass:u?.themeClass,onRender:u?.onRender,compitableColumn:f(e,[`columns`,`column`]),inlineThemeDisabled:i,mergedSize:o}},render(){let e=this.$slots.default,t=e?i(e()):[];t.length;let{contentClass:n,labelClass:r,compitableColumn:a,labelPlacement:s,labelAlign:c,mergedSize:l,bordered:d,title:f,cssVars:p,mergedClsPrefix:m,separator:g,onRender:_}=this;_?.();let v=t.filter(e=>S(e)),b=v.reduce((e,t,i)=>{let o=t.props||{},c=v.length-1===i,l=[`label`in o?o.label:y(t,`label`)],f=[y(t)],p=o.span||1,h=e.span;e.span+=p;let _=o.labelStyle||o[`label-style`]||this.labelStyle,b=o.contentStyle||o[`content-style`]||this.contentStyle;if(s===`left`)d?e.row.push(u(`th`,{class:[`${m}-descriptions-table-header`,r],colspan:1,style:_},l),u(`td`,{class:[`${m}-descriptions-table-content`,n],colspan:c?(a-h)*2+1:p*2-1,style:b},f)):e.row.push(u(`td`,{class:`${m}-descriptions-table-content`,colspan:c?(a-h)*2:p*2},u(`span`,{class:[`${m}-descriptions-table-content__label`,r],style:_},[...l,g&&u(`span`,{class:`${m}-descriptions-separator`},g)]),u(`span`,{class:[`${m}-descriptions-table-content__content`,n],style:b},f)));else{let t=c?(a-h)*2:p*2;e.row.push(u(`th`,{class:[`${m}-descriptions-table-header`,r],colspan:t,style:_},l)),e.secondRow.push(u(`td`,{class:[`${m}-descriptions-table-content`,n],colspan:t,style:b},f))}return(e.span>=a||c)&&(e.span=0,e.row.length&&(e.rows.push(e.row),e.row=[]),s!==`left`&&e.secondRow.length&&(e.rows.push(e.secondRow),e.secondRow=[])),e},{span:0,row:[],secondRow:[],rows:[]}).rows.map(e=>u(`tr`,{class:`${m}-descriptions-table-row`},e));return u(`div`,{style:p,class:[`${m}-descriptions`,this.themeClass,`${m}-descriptions--${s}-label-placement`,`${m}-descriptions--${c}-label-align`,`${m}-descriptions--${l}-size`,d&&`${m}-descriptions--bordered`]},f||this.$slots.header?u(`div`,{class:`${m}-descriptions-header`},f||o(this,`header`)):null,u(`div`,{class:`${m}-descriptions-table-wrapper`},u(`table`,{class:`${m}-descriptions-table`},u(`tbody`,null,s===`top`&&u(`tr`,{class:`${m}-descriptions-table-row`,style:{visibility:`collapse`}},h(a*2,u(`td`,null))),b))))}}),w={label:String,span:{type:Number,default:1},labelClass:String,labelStyle:[Object,String],contentClass:String,contentStyle:[Object,String]},T=m({name:`DescriptionsItem`,[x]:!0,props:w,slots:Object,render(){return null}});export{C as n,T as t};