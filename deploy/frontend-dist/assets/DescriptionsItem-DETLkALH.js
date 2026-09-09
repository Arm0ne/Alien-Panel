import{$i as e,Dr as t,Er as n,Ki as r,Kr as i,Qn as a,_i as o,bi as s,ea as c,gi as l,hi as u,hr as d,mi as f,mr as p,ni as m,pi as h,vi as g,yi as _}from"./router-C2u90Eon.js";import{a as v}from"./index-CRoWbJ3Q.js";function y(e,t=`default`,n=[]){let{children:r}=e;if(typeof r==`object`&&r&&!Array.isArray(r)){let e=r[t];if(typeof e==`function`)return e()}return n}var b=h([f(`descriptions`,{fontSize:`var(--n-font-size)`},[f(`descriptions-separator`,`
 display: inline-block;
 margin: 0 8px 0 2px;
 `),f(`descriptions-table-wrapper`,[f(`descriptions-table`,[f(`descriptions-table-row`,[f(`descriptions-table-header`,{padding:`var(--n-th-padding)`}),f(`descriptions-table-content`,{padding:`var(--n-td-padding)`})])])]),o(`bordered`,[f(`descriptions-table-wrapper`,[f(`descriptions-table`,[f(`descriptions-table-row`,[h(`&:last-child`,[f(`descriptions-table-content`,{paddingBottom:0})])])])])]),l(`left-label-placement`,[f(`descriptions-table-content`,[h(`> *`,{verticalAlign:`top`})])]),l(`left-label-align`,[h(`th`,{textAlign:`left`})]),l(`center-label-align`,[h(`th`,{textAlign:`center`})]),l(`right-label-align`,[h(`th`,{textAlign:`right`})]),l(`bordered`,[f(`descriptions-table-wrapper`,`
 border-radius: var(--n-border-radius);
 overflow: hidden;
 background: var(--n-merged-td-color);
 border: 1px solid var(--n-merged-border-color);
 `,[f(`descriptions-table`,[f(`descriptions-table-row`,[h(`&:not(:last-child)`,[f(`descriptions-table-content`,{borderBottom:`1px solid var(--n-merged-border-color)`}),f(`descriptions-table-header`,{borderBottom:`1px solid var(--n-merged-border-color)`})]),f(`descriptions-table-header`,`
 font-weight: 400;
 background-clip: padding-box;
 background-color: var(--n-merged-th-color);
 `,[h(`&:not(:last-child)`,{borderRight:`1px solid var(--n-merged-border-color)`})]),f(`descriptions-table-content`,[h(`&:not(:last-child)`,{borderRight:`1px solid var(--n-merged-border-color)`})])])])])]),f(`descriptions-header`,`
 font-weight: var(--n-th-font-weight);
 font-size: 18px;
 transition: color .3s var(--n-bezier);
 line-height: var(--n-line-height);
 margin-bottom: 16px;
 color: var(--n-title-text-color);
 `),f(`descriptions-table-wrapper`,`
 transition:
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `,[f(`descriptions-table`,`
 width: 100%;
 border-collapse: separate;
 border-spacing: 0;
 box-sizing: border-box;
 `,[f(`descriptions-table-row`,`
 box-sizing: border-box;
 transition: border-color .3s var(--n-bezier);
 `,[f(`descriptions-table-header`,`
 font-weight: var(--n-th-font-weight);
 line-height: var(--n-line-height);
 display: table-cell;
 box-sizing: border-box;
 color: var(--n-th-text-color);
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `),f(`descriptions-table-content`,`
 vertical-align: top;
 line-height: var(--n-line-height);
 display: table-cell;
 box-sizing: border-box;
 color: var(--n-td-text-color);
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `,[u(`content`,`
 transition: color .3s var(--n-bezier);
 display: inline-block;
 color: var(--n-td-text-color);
 `)]),u(`label`,`
 font-weight: var(--n-th-font-weight);
 transition: color .3s var(--n-bezier);
 display: inline-block;
 margin-right: 14px;
 color: var(--n-th-text-color);
 `)])])])]),f(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color);
 --n-merged-td-color: var(--n-td-color);
 --n-merged-border-color: var(--n-border-color);
 `),_(f(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color-modal);
 --n-merged-td-color: var(--n-td-color-modal);
 --n-merged-border-color: var(--n-border-color-modal);
 `)),s(f(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color-popover);
 --n-merged-td-color: var(--n-td-color-popover);
 --n-merged-border-color: var(--n-border-color-popover);
 `))]),x=`DESCRIPTION_ITEM_FLAG`;function S(e){return typeof e==`object`&&e&&!Array.isArray(e)?e.type&&e.type.DESCRIPTION_ITEM_FLAG:!1}var C=e({name:`Descriptions`,props:Object.assign(Object.assign({},a.props),{title:String,column:{type:Number,default:3},columns:Number,labelPlacement:{type:String,default:`top`},labelAlign:{type:String,default:`left`},separator:{type:String,default:`:`},size:String,bordered:Boolean,labelClass:String,labelStyle:[Object,String],contentClass:String,contentStyle:[Object,String]}),slots:Object,setup(e){let{mergedClsPrefixRef:t,inlineThemeDisabled:n,mergedComponentPropsRef:o}=d(e),s=r(()=>e.size||o?.value?.Descriptions?.size||`medium`),c=a(`Descriptions`,`-descriptions`,b,v,e,t),l=r(()=>{let{bordered:t}=e,n=s.value,{common:{cubicBezierEaseInOut:r},self:{titleTextColor:i,thColor:a,thColorModal:o,thColorPopover:l,thTextColor:u,thFontWeight:d,tdTextColor:f,tdColor:p,tdColorModal:m,tdColorPopover:h,borderColor:_,borderColorModal:v,borderColorPopover:y,borderRadius:b,lineHeight:x,[g(`fontSize`,n)]:S,[g(t?`thPaddingBordered`:`thPadding`,n)]:C,[g(t?`tdPaddingBordered`:`tdPadding`,n)]:w}}=c.value;return{"--n-title-text-color":i,"--n-th-padding":C,"--n-td-padding":w,"--n-font-size":S,"--n-bezier":r,"--n-th-font-weight":d,"--n-line-height":x,"--n-th-text-color":u,"--n-td-text-color":f,"--n-th-color":a,"--n-th-color-modal":o,"--n-th-color-popover":l,"--n-td-color":p,"--n-td-color-modal":m,"--n-td-color-popover":h,"--n-border-radius":b,"--n-border-color":_,"--n-border-color-modal":v,"--n-border-color-popover":y}}),u=n?p(`descriptions`,r(()=>{let t=``,{bordered:n}=e;return n&&(t+=`a`),t+=s.value[0],t}),l,e):void 0;return{mergedClsPrefix:t,cssVars:n?void 0:l,themeClass:u?.themeClass,onRender:u?.onRender,compitableColumn:i(e,[`columns`,`column`]),inlineThemeDisabled:n,mergedSize:s}},render(){let e=this.$slots.default,r=e?t(e()):[];r.length;let{contentClass:i,labelClass:a,compitableColumn:o,labelPlacement:s,labelAlign:l,mergedSize:u,bordered:d,title:f,cssVars:p,mergedClsPrefix:h,separator:g,onRender:_}=this;_?.();let v=r.filter(e=>S(e)),b=v.reduce((e,t,n)=>{let r=t.props||{},l=v.length-1===n,u=[`label`in r?r.label:y(t,`label`)],f=[y(t)],p=r.span||1,m=e.span;e.span+=p;let _=r.labelStyle||r[`label-style`]||this.labelStyle,b=r.contentStyle||r[`content-style`]||this.contentStyle;if(s===`left`)d?e.row.push(c(`th`,{class:[`${h}-descriptions-table-header`,a],colspan:1,style:_},u),c(`td`,{class:[`${h}-descriptions-table-content`,i],colspan:l?(o-m)*2+1:p*2-1,style:b},f)):e.row.push(c(`td`,{class:`${h}-descriptions-table-content`,colspan:l?(o-m)*2:p*2},c(`span`,{class:[`${h}-descriptions-table-content__label`,a],style:_},[...u,g&&c(`span`,{class:`${h}-descriptions-separator`},g)]),c(`span`,{class:[`${h}-descriptions-table-content__content`,i],style:b},f)));else{let t=l?(o-m)*2:p*2;e.row.push(c(`th`,{class:[`${h}-descriptions-table-header`,a],colspan:t,style:_},u)),e.secondRow.push(c(`td`,{class:[`${h}-descriptions-table-content`,i],colspan:t,style:b},f))}return(e.span>=o||l)&&(e.span=0,e.row.length&&(e.rows.push(e.row),e.row=[]),s!==`left`&&e.secondRow.length&&(e.rows.push(e.secondRow),e.secondRow=[])),e},{span:0,row:[],secondRow:[],rows:[]}).rows.map(e=>c(`tr`,{class:`${h}-descriptions-table-row`},e));return c(`div`,{style:p,class:[`${h}-descriptions`,this.themeClass,`${h}-descriptions--${s}-label-placement`,`${h}-descriptions--${l}-label-align`,`${h}-descriptions--${u}-size`,d&&`${h}-descriptions--bordered`]},f||this.$slots.header?c(`div`,{class:`${h}-descriptions-header`},f||n(this,`header`)):null,c(`div`,{class:`${h}-descriptions-table-wrapper`},c(`table`,{class:`${h}-descriptions-table`},c(`tbody`,null,s===`top`&&c(`tr`,{class:`${h}-descriptions-table-row`,style:{visibility:`collapse`}},m(o*2,c(`td`,null))),b))))}}),w={label:String,span:{type:Number,default:1},labelClass:String,labelStyle:[Object,String],contentClass:String,contentStyle:[Object,String]},T=e({name:`DescriptionsItem`,[x]:!0,props:w,slots:Object,render(){return null}});export{C as n,T as t};