import{Bi as e,Fr as t,Kr as n,Vn as r,Yi as i,ai as a,ar as o,ci as s,di as c,fi as l,gr as u,ir as d,li as f,oi as p,qi as m,si as h,ui as g}from"./router-DxH1-4bP.js";import{n as _}from"./Space-Ba7xxcgp.js";import{s as v}from"./index-BYJvH7k2.js";function y(e,t=`default`,n=[]){let{children:r}=e;if(typeof r==`object`&&r&&!Array.isArray(r)){let e=r[t];if(typeof e==`function`)return e()}return n}var b=a([p(`descriptions`,{fontSize:`var(--n-font-size)`},[p(`descriptions-separator`,`
 display: inline-block;
 margin: 0 8px 0 2px;
 `),p(`descriptions-table-wrapper`,[p(`descriptions-table`,[p(`descriptions-table-row`,[p(`descriptions-table-header`,{padding:`var(--n-th-padding)`}),p(`descriptions-table-content`,{padding:`var(--n-td-padding)`})])])]),f(`bordered`,[p(`descriptions-table-wrapper`,[p(`descriptions-table`,[p(`descriptions-table-row`,[a(`&:last-child`,[p(`descriptions-table-content`,{paddingBottom:0})])])])])]),s(`left-label-placement`,[p(`descriptions-table-content`,[a(`> *`,{verticalAlign:`top`})])]),s(`left-label-align`,[a(`th`,{textAlign:`left`})]),s(`center-label-align`,[a(`th`,{textAlign:`center`})]),s(`right-label-align`,[a(`th`,{textAlign:`right`})]),s(`bordered`,[p(`descriptions-table-wrapper`,`
 border-radius: var(--n-border-radius);
 overflow: hidden;
 background: var(--n-merged-td-color);
 border: 1px solid var(--n-merged-border-color);
 `,[p(`descriptions-table`,[p(`descriptions-table-row`,[a(`&:not(:last-child)`,[p(`descriptions-table-content`,{borderBottom:`1px solid var(--n-merged-border-color)`}),p(`descriptions-table-header`,{borderBottom:`1px solid var(--n-merged-border-color)`})]),p(`descriptions-table-header`,`
 font-weight: 400;
 background-clip: padding-box;
 background-color: var(--n-merged-th-color);
 `,[a(`&:not(:last-child)`,{borderRight:`1px solid var(--n-merged-border-color)`})]),p(`descriptions-table-content`,[a(`&:not(:last-child)`,{borderRight:`1px solid var(--n-merged-border-color)`})])])])])]),p(`descriptions-header`,`
 font-weight: var(--n-th-font-weight);
 font-size: 18px;
 transition: color .3s var(--n-bezier);
 line-height: var(--n-line-height);
 margin-bottom: 16px;
 color: var(--n-title-text-color);
 `),p(`descriptions-table-wrapper`,`
 transition:
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `,[p(`descriptions-table`,`
 width: 100%;
 border-collapse: separate;
 border-spacing: 0;
 box-sizing: border-box;
 `,[p(`descriptions-table-row`,`
 box-sizing: border-box;
 transition: border-color .3s var(--n-bezier);
 `,[p(`descriptions-table-header`,`
 font-weight: var(--n-th-font-weight);
 line-height: var(--n-line-height);
 display: table-cell;
 box-sizing: border-box;
 color: var(--n-th-text-color);
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `),p(`descriptions-table-content`,`
 vertical-align: top;
 line-height: var(--n-line-height);
 display: table-cell;
 box-sizing: border-box;
 color: var(--n-td-text-color);
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `,[h(`content`,`
 transition: color .3s var(--n-bezier);
 display: inline-block;
 color: var(--n-td-text-color);
 `)]),h(`label`,`
 font-weight: var(--n-th-font-weight);
 transition: color .3s var(--n-bezier);
 display: inline-block;
 margin-right: 14px;
 color: var(--n-th-text-color);
 `)])])])]),p(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color);
 --n-merged-td-color: var(--n-td-color);
 --n-merged-border-color: var(--n-border-color);
 `),c(p(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color-modal);
 --n-merged-td-color: var(--n-td-color-modal);
 --n-merged-border-color: var(--n-border-color-modal);
 `)),l(p(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color-popover);
 --n-merged-td-color: var(--n-td-color-popover);
 --n-merged-border-color: var(--n-border-color-popover);
 `))]),x=`DESCRIPTION_ITEM_FLAG`;function S(e){return typeof e==`object`&&e&&!Array.isArray(e)?e.type&&e.type.DESCRIPTION_ITEM_FLAG:!1}var C=m({name:`Descriptions`,props:Object.assign(Object.assign({},r.props),{title:String,column:{type:Number,default:3},columns:Number,labelPlacement:{type:String,default:`top`},labelAlign:{type:String,default:`left`},separator:{type:String,default:`:`},size:String,bordered:Boolean,labelClass:String,labelStyle:[Object,String],contentClass:String,contentStyle:[Object,String]}),slots:Object,setup(n){let{mergedClsPrefixRef:i,inlineThemeDisabled:a,mergedComponentPropsRef:s}=o(n),c=e(()=>n.size||s?.value?.Descriptions?.size||`medium`),l=r(`Descriptions`,`-descriptions`,b,v,n,i),u=e(()=>{let{bordered:e}=n,t=c.value,{common:{cubicBezierEaseInOut:r},self:{titleTextColor:i,thColor:a,thColorModal:o,thColorPopover:s,thTextColor:u,thFontWeight:d,tdTextColor:f,tdColor:p,tdColorModal:m,tdColorPopover:h,borderColor:_,borderColorModal:v,borderColorPopover:y,borderRadius:b,lineHeight:x,[g(`fontSize`,t)]:S,[g(e?`thPaddingBordered`:`thPadding`,t)]:C,[g(e?`tdPaddingBordered`:`tdPadding`,t)]:w}}=l.value;return{"--n-title-text-color":i,"--n-th-padding":C,"--n-td-padding":w,"--n-font-size":S,"--n-bezier":r,"--n-th-font-weight":d,"--n-line-height":x,"--n-th-text-color":u,"--n-td-text-color":f,"--n-th-color":a,"--n-th-color-modal":o,"--n-th-color-popover":s,"--n-td-color":p,"--n-td-color-modal":m,"--n-td-color-popover":h,"--n-border-radius":b,"--n-border-color":_,"--n-border-color-modal":v,"--n-border-color-popover":y}}),f=a?d(`descriptions`,e(()=>{let e=``,{bordered:t}=n;return t&&(e+=`a`),e+=c.value[0],e}),u,n):void 0;return{mergedClsPrefix:i,cssVars:a?void 0:u,themeClass:f?.themeClass,onRender:f?.onRender,compitableColumn:t(n,[`columns`,`column`]),inlineThemeDisabled:a,mergedSize:c}},render(){let e=this.$slots.default,t=e?u(e()):[];t.length;let{contentClass:r,labelClass:a,compitableColumn:o,labelPlacement:s,labelAlign:c,mergedSize:l,bordered:d,title:f,cssVars:p,mergedClsPrefix:m,separator:h,onRender:g}=this;g?.();let v=t.filter(e=>S(e)),b=v.reduce((e,t,n)=>{let c=t.props||{},l=v.length-1===n,u=[`label`in c?c.label:y(t,`label`)],f=[y(t)],p=c.span||1,g=e.span;e.span+=p;let _=c.labelStyle||c[`label-style`]||this.labelStyle,b=c.contentStyle||c[`content-style`]||this.contentStyle;if(s===`left`)d?e.row.push(i(`th`,{class:[`${m}-descriptions-table-header`,a],colspan:1,style:_},u),i(`td`,{class:[`${m}-descriptions-table-content`,r],colspan:l?(o-g)*2+1:p*2-1,style:b},f)):e.row.push(i(`td`,{class:`${m}-descriptions-table-content`,colspan:l?(o-g)*2:p*2},i(`span`,{class:[`${m}-descriptions-table-content__label`,a],style:_},[...u,h&&i(`span`,{class:`${m}-descriptions-separator`},h)]),i(`span`,{class:[`${m}-descriptions-table-content__content`,r],style:b},f)));else{let t=l?(o-g)*2:p*2;e.row.push(i(`th`,{class:[`${m}-descriptions-table-header`,a],colspan:t,style:_},u)),e.secondRow.push(i(`td`,{class:[`${m}-descriptions-table-content`,r],colspan:t,style:b},f))}return(e.span>=o||l)&&(e.span=0,e.row.length&&(e.rows.push(e.row),e.row=[]),s!==`left`&&e.secondRow.length&&(e.rows.push(e.secondRow),e.secondRow=[])),e},{span:0,row:[],secondRow:[],rows:[]}).rows.map(e=>i(`tr`,{class:`${m}-descriptions-table-row`},e));return i(`div`,{style:p,class:[`${m}-descriptions`,this.themeClass,`${m}-descriptions--${s}-label-placement`,`${m}-descriptions--${c}-label-align`,`${m}-descriptions--${l}-size`,d&&`${m}-descriptions--bordered`]},f||this.$slots.header?i(`div`,{class:`${m}-descriptions-header`},f||_(this,`header`)):null,i(`div`,{class:`${m}-descriptions-table-wrapper`},i(`table`,{class:`${m}-descriptions-table`},i(`tbody`,null,s===`top`&&i(`tr`,{class:`${m}-descriptions-table-row`,style:{visibility:`collapse`}},n(o*2,i(`td`,null))),b))))}}),w={label:String,span:{type:Number,default:1},labelClass:String,labelStyle:[Object,String],contentClass:String,contentStyle:[Object,String]},T=m({name:`DescriptionsItem`,[x]:!0,props:w,slots:Object,render(){return null}});export{C as n,T as t};